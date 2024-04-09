using Microsoft.EntityFrameworkCore;
using KitchenKomplete.Helpers;
using Microsoft.IdentityModel.Tokens;

namespace Models.Contexts;

public class KitchenKompleteContext : DbContext
{
    IConfiguration _configuration;
    public KitchenKompleteContext(DbContextOptions options, IConfiguration configuration) : base(options)
    {
        _configuration = configuration;
    }

    protected override void OnConfiguring(DbContextOptionsBuilder optionsBuilder)
    {
        if (System.Diagnostics.Debugger.IsAttached == false)
        {
            System.Diagnostics.Debugger.Launch();
        }
        var connString = _configuration.GetConnectionString("DB");
        if(connString.IsNullOrEmpty()){throw new Exception(message: "Connection string blank");}
        optionsBuilder.UseNpgsql(connString);
    }

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.Entity<FoodItem>().OwnsOne(item=>item.NutritionInfo, builder =>{builder.ToJson();});
    }


    public DbSet<Pantry> Pantries { get; set; } = null!;

    public DbSet<Recipe> Recipes { get; set; } = null!;

    public DbSet<FoodItem> FoodItem { get; set; } = null!;


}