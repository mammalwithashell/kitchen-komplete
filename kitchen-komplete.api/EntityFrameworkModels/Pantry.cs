namespace Models
{

    public class Pantry
    {

        public Guid Id {get; set;}
        public Guid OwnerId {get; set;}
        public List<Guid> Inventory {get; set;}
    }
}