using KitchenKomplete.Helpers;
using KitchenKomplete.Services;
using Microsoft.AspNetCore.Mvc;
using Microsoft.IdentityModel.Tokens;
using Models;

namespace KitchenKomplete.Controllers
{
    [Route("api/Recipes")]
    [ApiController]
    public class RecipesController : ControllerBase
    {
        private readonly ILogger<RecipesController> _logger;
        private readonly IRecipeService _recipeService;

        public RecipesController(ILogger<RecipesController> logger, IRecipeService recipeService)
        {
            _logger = logger;
            _recipeService = recipeService;
        }

        // GET: api/Recipes
        [HttpGet]
        public async Task<ActionResult<IEnumerable<Recipe>>> GetRecipes()
        {
            var response = await _recipeService.GetRecipesAsync();
            if (!response.ErrorMessages.IsNullOrEmpty() && !response.Data.IsNullOrEmpty()){
                return NoContent();
            }
            return Ok(response.Data);
        }

        // GET: api/Recipes/5
        [HttpGet("{id}")]
        public async Task<ActionResult<Recipe>> GetRecipe(Guid id)
        {
            var response = await _recipeService.GetRecipeAsync(id);
            if (response.Data != null) {
                throw new Exception();
            }
            return Ok(response.Data);
        }

        // PUT: api/Recipes/5
        // To protect from overposting attacks, see https://go.microsoft.com/fwlink/?linkid=2123754
        [HttpPut("{id}")]
        public async Task<IActionResult> PutRecipe(Guid id, Recipe recipe)
        {
            var response = await _recipeService.PutRecipe(id, recipe);
            if(response.ErrorMessages is not null){
                return HandleServiceResponse(response);
            }
            return Ok(recipe);
        }

        // POST: api/Recipes
        // To protect from overposting attacks, see https://go.microsoft.com/fwlink/?linkid=2123754
        [HttpPost]
        public async Task<IActionResult> PostRecipe(Recipe recipe)
        {
            var response = await _recipeService.PostRecipe(recipe);
            if (response.ErrorMessages is not null){
                return HandleServiceResponse(response);
            }
            return CreatedAtAction(nameof(GetRecipe), new { id = recipe.Id }, recipe);
        }

        // DELETE: api/Recipes/5
        [HttpDelete("{id}")]
        public async Task<IActionResult> DeleteRecipe(Guid id)
        {
            var response = await _recipeService.DeleteRecipe(id);
            if (response.ErrorMessages is not null){
                return HandleServiceResponse(response);
            }
            return NoContent();
        }

        private IActionResult HandleServiceResponse<T>(ServiceResponse<T> response) => response.ErrorMessages switch {
            [ServiceResponseHandler.BAD_REQUEST, ..] => BadRequest(..),
            [ServiceResponseHandler.NOT_FOUND, ..] => NotFound(new {
                errors = ..,

            }),
            _ => Ok(), 
        };

    }
}
