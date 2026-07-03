using System;
using System.Net.Http;
using System.Text;
using System.Threading.Tasks;
using Newtonsoft.Json;

class Program
{
    static readonly HttpClient client = new HttpClient();

    const string BASE_URL = "http://localhost:8080";
    const string API_KEY = "your-api-key";
    const string APP_NAME = "mytool";
    const string HWID = "test-hwid-123";

    static async Task<string> CreateKey()
    {
        var body = new
        {
            app = APP_NAME,
            duration = "30d"
        };

        var request = new HttpRequestMessage(HttpMethod.Post, $"{BASE_URL}/create-key");
        request.Headers.Add("X-Api-Key", API_KEY);
        request.Content = new StringContent(
            JsonConvert.SerializeObject(body),
            Encoding.UTF8,
            "application/json"
        );

        var response = await client.SendAsync(request);
        var json = await response.Content.ReadAsStringAsync();

        Console.WriteLine("Create Key:");
        Console.WriteLine(json);

        dynamic obj = JsonConvert.DeserializeObject(json);
        return obj.key;
    }

    static async Task Login(string key, string hwid)
    {
        var body = new
        {
            app = APP_NAME,
            key = key,
            hwid = hwid
        };

        var response = await client.PostAsync(
            $"{BASE_URL}/login",
            new StringContent(
                JsonConvert.SerializeObject(body),
                Encoding.UTF8,
                "application/json"
            )
        );

        Console.WriteLine(await response.Content.ReadAsStringAsync());
    }

    static async Task DeleteKey(string key)
    {
        var body = new
        {
            app = APP_NAME,
            key = key
        };

        var request = new HttpRequestMessage(HttpMethod.Post, $"{BASE_URL}/delete-key");
        request.Headers.Add("X-Api-Key", API_KEY);
        request.Content = new StringContent(
            JsonConvert.SerializeObject(body),
            Encoding.UTF8,
            "application/json"
        );

        var response = await client.SendAsync(request);

        Console.WriteLine(await response.Content.ReadAsStringAsync());
    }

    static async Task Main()
    {
        var key = await CreateKey();

        await Login(key, HWID);
        await Login(key, HWID);
        await Login(key, "another-hwid");

        await DeleteKey(key);
    }
}