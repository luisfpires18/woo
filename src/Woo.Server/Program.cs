using System.Reflection;
using System.Text.Json;

var builder = WebApplication.CreateBuilder(args);
var app = builder.Build();

// Only files the Vite build manifest lists (content-hashed output) may be cached forever.
// Everything else (index.html and any unversioned file, even under /assets/) must be revalidated.
var hashedFiles = ViteManifest.HashedFiles(app.Environment.WebRootPath);
var staticFiles = new StaticFileOptions
{
    OnPrepareResponse = ctx => ctx.Context.Response.Headers.CacheControl =
        hashedFiles.Contains(ctx.Context.Request.Path.Value ?? "")
            ? "public, max-age=31536000, immutable"
            : "no-cache",
};

// Published builds serve the React app from wwwroot.
app.UseStaticFiles(staticFiles);

app.MapGet("/api/health", () => Results.Ok(new { status = "ok" }));

var build = BuildIdentity.FromAssembly(typeof(Program).Assembly);
app.MapGet("/api/build", (HttpContext context) =>
{
    context.Response.Headers.CacheControl = "no-store";
    return Results.Ok(build);
});

// Unknown API routes stay API errors; they must never fall through to the SPA's index.html.
app.Map("/api/{**path}", () => Results.Problem(statusCode: StatusCodes.Status404NotFound));

// Client-side routes (paths without a file extension) get the SPA entry page.
app.MapFallbackToFile("index.html", staticFiles);

app.Run();

public partial class Program;

/// <summary>Reads the web root's Vite manifest (build.manifest) to find its content-hashed output files.</summary>
public static class ViteManifest
{
    private sealed record Chunk(string? File, string[]? Css, string[]? Assets);

    /// <returns>Request paths such as "/assets/index-AbCd1234.js"; empty when there is no manifest.</returns>
    public static HashSet<string> HashedFiles(string? webRoot)
    {
        var path = webRoot is null ? null : Path.Combine(webRoot, ".vite", "manifest.json");
        if (path is null || !File.Exists(path)) return [];
        var chunks = JsonSerializer.Deserialize<Dictionary<string, Chunk>>(File.ReadAllText(path), JsonSerializerOptions.Web) ?? [];
        return chunks.Values
            .SelectMany(c => (c.File is null ? [] : new[] { c.File }).Concat(c.Css ?? []).Concat(c.Assets ?? []))
            .Select(file => "/" + file)
            .ToHashSet(StringComparer.Ordinal);
    }
}

/// <summary>Build identity embedded at build time; see docs/technical/stack.md.</summary>
public sealed record BuildIdentity(string Version, string Commit, string BuiltAt, string Environment)
{
    public static BuildIdentity FromAssembly(Assembly assembly)
    {
        var values = assembly.GetCustomAttributes<AssemblyMetadataAttribute>()
            .ToDictionary(a => a.Key, a => a.Value ?? "");
        string Get(string key) => values.GetValueOrDefault(key, "");
        return new(Get("WooVersion"), Get("WooCommit"), Get("WooBuiltAt"), Get("WooEnvironment"));
    }
}
