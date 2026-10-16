using System.Net;
using System.Net.Http.Json;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Mvc.Testing;

namespace Woo.Server.Tests;

public sealed class HostTests : IClassFixture<HostTests.SpaFactory>
{
    private readonly HttpClient _client;

    public HostTests(SpaFactory factory) => _client = factory.CreateClient();

    private static CancellationToken Ct => TestContext.Current.CancellationToken;

    [Fact]
    public async Task Health_returns_ok_json()
    {
        var response = await _client.GetAsync("/api/health", Ct);

        Assert.Equal(HttpStatusCode.OK, response.StatusCode);
        var body = await response.Content.ReadFromJsonAsync<HealthBody>(Ct);
        Assert.Equal("ok", body?.Status);
    }

    [Fact]
    public async Task Build_identity_uses_VERSION_and_honest_local_fallback()
    {
        var response = await _client.GetAsync("/api/build", Ct);

        Assert.Equal(HttpStatusCode.OK, response.StatusCode);
        Assert.Equal("no-store", response.Headers.CacheControl?.ToString());
        var body = await response.Content.ReadFromJsonAsync<BuildIdentity>(Ct);
        // Tests run on a standalone build: version from VERSION, environment local, nothing invented.
        Assert.Equal(new BuildIdentity(ReadVersionFile(), "", "", "local"), body);
    }

    [Theory]
    [InlineData("/api/does-not-exist")]
    [InlineData("/api/health/extra")]
    [InlineData("/api/build/extra")]
    public async Task Unknown_api_route_is_404_problem_not_spa_html(string path)
    {
        var response = await _client.GetAsync(path, Ct);

        Assert.Equal(HttpStatusCode.NotFound, response.StatusCode);
        Assert.Equal("application/problem+json", response.Content.Headers.ContentType?.MediaType);
    }

    [Theory]
    [InlineData("/some/client/route")]
    [InlineData("/worlds")] // main page behind the title screen (task 011) loads directly
    [InlineData("/game/settlement")] // game-shell routes (task 012) load directly
    [InlineData("/index.html")]
    public async Task Entry_html_is_served_and_must_revalidate(string path)
    {
        var response = await _client.GetAsync(path, Ct);

        Assert.Equal(HttpStatusCode.OK, response.StatusCode);
        Assert.Equal("text/html", response.Content.Headers.ContentType?.MediaType);
        Assert.Contains(SpaFactory.Marker, await response.Content.ReadAsStringAsync(Ct));
        Assert.Equal("no-cache", response.Headers.CacheControl?.ToString());
    }

    [Fact]
    public async Task Hashed_assets_are_immutable()
    {
        var response = await _client.GetAsync("/assets/index-AbCd1234.js", Ct);

        Assert.Equal(HttpStatusCode.OK, response.StatusCode);
        Assert.Equal("public, max-age=31536000, immutable", response.Headers.CacheControl?.ToString());
    }

    [Theory]
    [InlineData("/robots.txt")]
    [InlineData("/assets/settings.json")] // under /assets/ but not in the build manifest
    [InlineData("/.vite/manifest.json")] // public anyway: it only names files index.html already references
    public async Task Unversioned_files_must_revalidate(string path)
    {
        var response = await _client.GetAsync(path, Ct);

        Assert.Equal(HttpStatusCode.OK, response.StatusCode);
        Assert.Equal("no-cache", response.Headers.CacheControl?.ToString());
    }

    private static string ReadVersionFile()
    {
        for (var dir = new DirectoryInfo(AppContext.BaseDirectory); dir is not null; dir = dir.Parent)
        {
            var path = Path.Combine(dir.FullName, "VERSION");
            if (File.Exists(path)) return File.ReadAllText(path).Trim();
        }
        throw new FileNotFoundException("VERSION not found above the test output folder.");
    }

    private sealed record HealthBody(string Status);

    // Stands in for the published wwwroot: entry page, Vite manifest, one hashed asset and unversioned files.
    public sealed class SpaFactory : WebApplicationFactory<Program>
    {
        public const string Marker = "<!-- woo-test-index -->";
        private readonly string _webRoot = Directory.CreateTempSubdirectory("woo-webroot-").FullName;

        public SpaFactory()
        {
            File.WriteAllText(Path.Combine(_webRoot, "index.html"), Marker);
            File.WriteAllText(Path.Combine(_webRoot, "robots.txt"), "User-agent: *");
            Directory.CreateDirectory(Path.Combine(_webRoot, "assets"));
            File.WriteAllText(Path.Combine(_webRoot, "assets", "index-AbCd1234.js"), "// hashed");
            File.WriteAllText(Path.Combine(_webRoot, "assets", "settings.json"), "{}");
            Directory.CreateDirectory(Path.Combine(_webRoot, ".vite"));
            File.WriteAllText(Path.Combine(_webRoot, ".vite", "manifest.json"),
                """{"index.html":{"file":"assets/index-AbCd1234.js","name":"index","src":"index.html","isEntry":true}}""");
        }

        protected override void ConfigureWebHost(IWebHostBuilder builder) => builder.UseWebRoot(_webRoot);

        protected override void Dispose(bool disposing)
        {
            base.Dispose(disposing);
            Directory.Delete(_webRoot, recursive: true);
        }
    }
}
