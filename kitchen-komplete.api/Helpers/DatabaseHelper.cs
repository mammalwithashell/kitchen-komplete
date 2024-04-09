
namespace KitchenKomplete.Helpers
{
    public static class DatabaseHelper{
        public static string? GetRDSConnectionString(IConfiguration configuration){

            string? dbname = configuration.GetConnectionString("DB");

            if (string.IsNullOrEmpty(dbname)) return null;

            string? username = Environment.GetEnvironmentVariable("ASPNETCORE_RDS_USERNAME");
            if (string.IsNullOrEmpty(username)) return null;

            string? password = Environment.GetEnvironmentVariable("ASPNETCORE_RDS_PASSWORD");
            if (string.IsNullOrEmpty(password)) return null;

            string? hostname = Environment.GetEnvironmentVariable("ASPNETCORE_RDS_HOSTNAME");
            if (string.IsNullOrEmpty(hostname)) return null;

            string? port = Environment.GetEnvironmentVariable("ASPNETCORE_RDS_PORT");
            if (string.IsNullOrEmpty(port)) return null;
            Console.WriteLine($"Host={hostname};Port={port};Username={username};Password={password};Database={dbname}");
            return $"Host={hostname};Port={port};Username={username};Password={password};Database={dbname}";
        }
    }
}