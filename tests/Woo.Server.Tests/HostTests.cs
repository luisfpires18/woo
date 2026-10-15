using System.Net;
using System.Net.Http.Json;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Mvc.Testing;

namespace Woo.Server.Tests;

public sealed class HostTests : IClassFixture<HostTests.SpaFactory>
{
    private readonly HttpClient _client;

    public HostTests(SpaFactory factory) => _client = factory.CreateClient();

    [Fact]
    public async Task Health_returns_ok_json()
    {
        var response = await _client.GetAsync("/api/health", TestContext.Current.CancellationToken);

        Assert.Equal(HttpStatusCode.OK, response.StatusCode);
        var body = await response.Content.ReadFromJsonAsync<HealthBody>(TestContext.Current.CancellationToken);
        Assert.Equal("ok", body?.Status);
    }

    [Theory]
    [InlineData("/api/does-not-exist")]
    [InlineData("/api/health/extra")]
    public async Task Unknown_api_route_is_404_problem_not_spa_html(string path)
    {
        var response = await _client.GetAsync(path, TestContext.Current.CancellationToken);

        Assert.Equal(HttpStatusCode.NotFound, response.StatusCode);
        Assert.Equal("application/problem+json", response.Content.Headers.ContentType?.MediaType);
    }

    [Fact]
    public async Task Client_route_falls_back_to_index_html()
    {
        var response = await _client.GetAsync("/some/client/route", TestContext.Current.CancellationToken);

        Assert.Equal(HttpStatusCode.OK, response.StatusCode);
        Assert.Equal("text/html", response.Content.Headers.ContentType?.MediaType);
        Assert.Contains(SpaFactory.Marker, await response.Content.ReadAsStringAsync(TestContext.Current.CancellationToken));
    }

    private sealed record HealthBody(string Status);

    // Stands in for the published wwwroot with a temporary index.html.
    public sealed class SpaFactory : WebApplicationFactory<Program>
    {
        public const string Marker = "<!-- woo-test-index -->";
        private readonly string _webRoot = Directory.CreateTempSubdirectory("woo-webroot-").FullName;

        public SpaFactory() => File.WriteAllText(Path.Combine(_webRoot, "index.html"), Marker);

        protected override void ConfigureWebHost(IWebHostBuilder builder) => builder.UseWebRoot(_webRoot);

        protected override void Dispose(bool disposing)
        {
            base.Dispose(disposing);
            Directory.Delete(_webRoot, recursive: true);
        }
    }
}
