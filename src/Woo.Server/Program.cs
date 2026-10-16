using System.Reflection;

var builder = WebApplication.CreateBuilder(args);
var app = builder.Build();

// Vite emits content-hashed files only under /assets/, so those may be cached forever.
// Everything else (index.html and any unversioned file) must be revalidated.
var staticFiles = new StaticFileOptions
{
    OnPrepareResponse = ctx => ctx.Context.Response.Headers.CacheControl =
        ctx.Context.Request.Path.StartsWithSegments("/assets")
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
