#include <iostream>
#include <cpr/cpr.h>
#include <nlohmann/json.hpp>

using json = nlohmann::json;

const std::string BASE_URL = "http://localhost:8080";
const std::string API_KEY = "your-api-key";
const std::string APP_NAME = "mytool";
const std::string HWID = "test-hwid-123";

std::string createKey()
{
    json body = {
        {"app", APP_NAME},
        {"duration", "30d"}
    };

    auto r = cpr::Post(
        cpr::Url{BASE_URL + "/create-key"},
        cpr::Header{
            {"Content-Type","application/json"},
            {"X-Api-Key",API_KEY}
        },
        cpr::Body{body.dump()}
    );

    std::cout << r.text << std::endl;

    return json::parse(r.text)["key"];
}

void login(std::string key, std::string hwid)
{
    json body = {
        {"app", APP_NAME},
        {"key", key},
        {"hwid", hwid}
    };

    auto r = cpr::Post(
        cpr::Url{BASE_URL + "/login"},
        cpr::Header{
            {"Content-Type","application/json"}
        },
        cpr::Body{body.dump()}
    );

    std::cout << r.text << std::endl;
}

void deleteKey(std::string key)
{
    json body = {
        {"app", APP_NAME},
        {"key", key}
    };

    auto r = cpr::Post(
        cpr::Url{BASE_URL + "/delete-key"},
        cpr::Header{
            {"Content-Type","application/json"},
            {"X-Api-Key",API_KEY}
        },
        cpr::Body{body.dump()}
    );

    std::cout << r.text << std::endl;
}

int main()
{
    std::string key = createKey();

    login(key, HWID);
    login(key, HWID);
    login(key, "another-hwid");

    deleteKey(key);

    return 0;
}