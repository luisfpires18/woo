var builder = WebApplication.CreateBuilder(args);
var app = builder.Build();

// Published builds serve the React app from wwwroot.
app.UseStaticFiles();

app.MapGet("/api/health", () => Results.Ok(new { status = "ok" }));

// Unknown API routes stay API errors; they must never fall through to the SPA's index.html.
app.Map("/api/{**path}", () => Results.Problem(statusCode: StatusCodes.Status404NotFound));

// Client-side routes (paths without a file extension) get the SPA entry page.
app.MapFallbackToFile("index.html");

app.Run();

public partial class Program;
