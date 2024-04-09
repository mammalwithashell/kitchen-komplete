using Microsoft.EntityFrameworkCore;
using Models;
using Models.Contexts;
using KitchenKomplete.Helpers;

namespace KitchenKomplete.Services
{

    public class RecipeService
    : IRecipeService
    {

        private readonly IServiceScopeFactory _scopeFactory;
        public RecipeService(IServiceScopeFactory scopeFactory)
        {
            _scopeFactory = scopeFactory;
        }

        public async Task<ServiceResponse<List<Recipe>>> GetRecipesAsync()
        {
            ServiceResponse<List<Recipe>> response = new();
            using (var scope = _scopeFactory.CreateScope())
            {
                var db = scope.ServiceProvider.GetService<KitchenKompleteContext>();
                response = new ServiceResponse<List<Recipe>>
                {
                    Data = await db.Recipes.ToListAsync()
                };
            }

            return response;
        }

        public async Task<ServiceResponse<Recipe>> GetRecipeAsync(Guid id)
        {
            var response = new ServiceResponse<Recipe>();
            var recipe = new Recipe();
            using (var scope = _scopeFactory.CreateScope())
            {
                var db = scope.ServiceProvider.GetService<KitchenKompleteContext>();
                recipe = await db.Recipes.FindAsync(id);
            }

            if (recipe == null)
            {
                response.ErrorMessages = [ServiceResponseHandler.BAD_REQUEST, $"No matching recipes w/ Id {id}"];
                return response;
            }

            response.Data = recipe;
            return response;
        }

        public async Task<ServiceResponse<Recipe>> PutRecipe(Guid id, Recipe recipe)
        {
            var response = new ServiceResponse<Recipe>();
            if (id != recipe.Id)
            {
                response.ErrorMessages = [ServiceResponseHandler.BAD_REQUEST];
                return response;
            }
            using (var scope = _scopeFactory.CreateScope())
            {
                var db = scope.ServiceProvider.GetService<KitchenKompleteContext>();
                db.Entry(recipe).State = EntityState.Modified;
                try
                {
                    await db.SaveChangesAsync();
                }
                catch (DbUpdateConcurrencyException)
                {
                    if (!RecipeExists(id, db))
                    {
                        response.ErrorMessages = [ServiceResponseHandler.NOT_FOUND];
                        return response;
                    }
                    else
                    {
                        throw;
                    }
                }
            }


            response.Data = recipe;
            return response;
        }

        public async Task<ServiceResponse<Recipe>> PostRecipe(Recipe recipe)
        {
            var response = new ServiceResponse<Recipe>();
            using (var scope = _scopeFactory.CreateScope())
            {
                var db = scope.ServiceProvider.GetService<KitchenKompleteContext>();
                db.Recipes.Add(recipe);
                response.Data = recipe;
                await db.SaveChangesAsync();
            }


            return response;
        }

        public async Task<ServiceResponse<Recipe>> DeleteRecipe(Guid id)
        {
            var response = new ServiceResponse<Recipe>();

            var recipe = new Recipe();
            using (var scope = _scopeFactory.CreateScope())
            {
                var db = scope.ServiceProvider.GetService<KitchenKompleteContext>();

                recipe = await db.Recipes.FindAsync(id);
                if (recipe is null)
                {
                    response.ErrorMessages = [ServiceResponseHandler.NOT_FOUND];
                    return response;
                }

                db.Recipes.Remove(recipe);
                await db.SaveChangesAsync();
            }
            response.Data = recipe;
            return response;
        }

        private static bool RecipeExists(Guid id, KitchenKompleteContext? db)
        {
            return db.Recipes.Any(e => e.Id == id);
        }
    }

    public interface IRecipeService
    {
        Task<ServiceResponse<List<Recipe>>> GetRecipesAsync();

        Task<ServiceResponse<Recipe>> GetRecipeAsync(Guid id);

        Task<ServiceResponse<Recipe>> PutRecipe(Guid id, Recipe recipe);

        Task<ServiceResponse<Recipe>> PostRecipe(Recipe recipe);

        Task<ServiceResponse<Recipe>> DeleteRecipe(Guid id);
    }
}
