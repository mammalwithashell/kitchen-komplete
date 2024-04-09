using Microsoft.AspNetCore.Mvc;

namespace Models{
    public class ServiceResponse<T>{
        public List<string>? ErrorMessages {get; set;}

        public IActionResult? ControllerResponse {get; set;}
        public T? Data {get; set;}
    }
}