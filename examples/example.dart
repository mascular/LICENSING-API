import 'dart:convert';
import 'package:http/http.dart' as http;

const baseUrl = "http://localhost:8080";
const apiKey = "your-api-key";
const appName = "mytool";
const hwid = "test-hwid-123";

Future<String> createKey() async {
  final response = await http.post(
    Uri.parse("$baseUrl/create-key"),
    headers: {
      "Content-Type": "application/json",
      "X-Api-Key": apiKey,
    },
    body: jsonEncode({
      "app": appName,
      "duration": "30d",
    }),
  );

  print(response.body);

  return jsonDecode(response.body)["key"];
}

Future<void> login(String key, String hwid) async {
  final response = await http.post(
    Uri.parse("$baseUrl/login"),
    headers: {
      "Content-Type": "application/json",
    },
    body: jsonEncode({
      "app": appName,
      "key": key,
      "hwid": hwid,
    }),
  );

  print(response.body);
}

Future<void> deleteKey(String key) async {
  final response = await http.post(
    Uri.parse("$baseUrl/delete-key"),
    headers: {
      "Content-Type": "application/json",
      "X-Api-Key": apiKey,
    },
    body: jsonEncode({
      "app": appName,
      "key": key,
    }),
  );

  print(response.body);
}

Future<void> main() async {
  final key = await createKey();

  await login(key, hwid);
  await login(key, hwid);
  await login(key, "another-hwid");

  await deleteKey(key);
}