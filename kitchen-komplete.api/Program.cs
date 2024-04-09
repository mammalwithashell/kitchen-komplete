using DotNetEnv;
using KitchenKomplete.Services;
using Models.Contexts;

var builder = WebApplication.CreateBuilder(args);


// Add services to the container.
ConfigureServices(builder.Services);
builder.Services.AddControllers();
// Learn more about configuring Swagger/OpenAPI at https://aka.ms/aspnetcore/swashbuckle
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();
builder.Configuration.AddJsonFile("appsettings.json", optional: false, reloadOnChange: true);
var app = builder.Build();


// Configure the HTTP request pipeline.
if (app.Environment.IsDevelopment())
{
    Env.Load();
    app.UseSwagger();
    app.UseSwaggerUI();
}
// using (var scope = app.Services.CreateScope())
// {
//     var services = scope.ServiceProvider;

//     var context = services.GetRequiredService<KitchenKompleteContext>();

//     context.Database.EnsureCreated();
//     // DbInitializer.Initialize(context);
// }

app.UseHttpsRedirection();

app.UseAuthorization();


app.MapControllers();

app.Run();

static void ConfigureServices(IServiceCollection services){
    services.AddDbContext<KitchenKompleteContext>();
    services.AddSingleton<IRecipeService, RecipeService>();
}