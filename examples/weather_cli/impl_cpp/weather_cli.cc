// Reference C++ implementation of weather-cli.
//
// This file is the *legacy* artifact in the worked example: the
// archaeologist distills it into ../system.dx, the architect refines that
// spec, and the implementer regenerates the system in another language
// (see ../impl_python/) without ever reading this file.
//
// Network calls and JSON parsing are mocked for brevity. The structural
// contract surface (env-var handling, cache TTL by mtime, exit codes,
// stdout/stderr split) is intact and matches the contracts in
// ../system.dx.

#include <chrono>
#include <cstdlib>
#include <fstream>
#include <iostream>
#include <string>
#include <sys/stat.h>

namespace {

constexpr int kCacheTtlSeconds = 600;

std::string CacheFilePath() {
  const char* home = std::getenv("HOME");
  if (home == nullptr) {
    // Fall back to the working directory rather than crashing; the
    // contracts do not pin a path.
    return ".weather_cache.json";
  }
  return std::string(home) + "/.weather_cache.json";
}

bool IsCacheValid(const std::string& path) {
  struct stat result;
  if (stat(path.c_str(), &result) != 0) {
    return false;
  }
  const auto now = std::chrono::system_clock::now();
  const auto cache_time =
      std::chrono::system_clock::from_time_t(result.st_mtime);
  const auto duration =
      std::chrono::duration_cast<std::chrono::seconds>(now - cache_time);
  return duration.count() < kCacheTtlSeconds;
}

// Simulated upstream call. A real implementation would issue an HTTP
// request to the OpenMeteo (or equivalent) endpoint.
std::string FetchFromUpstream(const std::string& zip,
                              const std::string& /*api_key*/) {
  return std::string("{\"zip\": \"") + zip +
         "\", \"temp\": \"72F\", \"condition\": \"Sunny\"}";
}

}  // namespace

int main(int argc, char* argv[]) {
  if (argc < 2) {
    std::cerr << "Usage: weather_cli <zipcode> [--json]\n";
    return 1;
  }

  const std::string zip = argv[1];
  const bool json_output =
      (argc == 3 && std::string(argv[2]) == "--json");

  const char* api_key_env = std::getenv("WEATHER_API_KEY");
  if (api_key_env == nullptr || api_key_env[0] == '\0') {
    std::cerr << "Error: WEATHER_API_KEY environment variable not set.\n";
    return 1;
  }

  const std::string cache_path = CacheFilePath();
  std::string weather_data;

  if (IsCacheValid(cache_path)) {
    std::ifstream cache(cache_path);
    std::getline(cache, weather_data);
  }

  if (weather_data.empty()) {
    weather_data = FetchFromUpstream(zip, api_key_env);
    std::ofstream cache(cache_path);
    cache << weather_data;
  }

  if (json_output) {
    std::cout << weather_data << "\n";
  } else {
    // Human-readable summary; exact phrasing is unconstrained.
    std::cout << "Weather for " << zip << ": 72F, Sunny\n";
  }

  return 0;
}
