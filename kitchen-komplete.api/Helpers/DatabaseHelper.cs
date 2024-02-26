
using System;
namespace KitchenKomplete.Helpers
{
    public class DatabaseHelper{
        public static string? GetRDSConnectionString(){

            string? dbname = Environment.GetEnvironmentVariable("RDS_DB_NAME");

            if (string.IsNullOrEmpty(dbname)) return null;

            string? username = Environment.GetEnvironmentVariable("RDS_USERNAME");
            if (string.IsNullOrEmpty(username)) return null;

            string? password = Environment.GetEnvironmentVariable("RDS_PASSWORD");
            if (string.IsNullOrEmpty(password)) return null;

            string? hostname = Environment.GetEnvironmentVariable("RDS_HOSTNAME");
            if (string.IsNullOrEmpty(hostname)) return null;

            string? port = Environment.GetEnvironmentVariable("RDS_PORT");
            if (string.IsNullOrEmpty(port)) return null;


            return $"Host={hostname};Port={port};;User ID={username};Password={password};Database={dbname}";
        }
    }
}