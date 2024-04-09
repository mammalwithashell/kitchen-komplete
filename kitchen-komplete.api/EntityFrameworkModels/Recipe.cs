namespace Models;

public class Recipe
{
    public Guid Id { get; set; }
    public string Name { get; set; }
    public List<Guid> Ingredients {get; set;}
    public DateTime LastUpdate { get; set;}
    public DateTime InsertDate { get; set;}
}