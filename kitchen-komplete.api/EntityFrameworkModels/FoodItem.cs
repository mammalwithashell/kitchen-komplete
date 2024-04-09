using System.ComponentModel.DataAnnotations.Schema;
using Microsoft.EntityFrameworkCore;

namespace Models;

public class FoodItem
{
    public Guid Id { get; set; }
    public string Name {get; set;}
    public NutritionInfo NutritionInfo { get; set; }

}

public class NutritionInfo
{
    public int GramsProtien { get; set; }
    public int GramsCarbohydrates { get; set; }

    public int GramsFat { get; set; }
}