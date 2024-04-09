namespace KitchenKomplete.Test;
using Microsoft.VisualStudio.TestPlatform.ObjectModel.Client;
using Xunit;
using Models.Contexts;
using Moq;

public class DatabaseTest
{
    private ITestRunConfiguration _configuration;
    public DatabaseTest (ITestRunConfiguration configuration){
        _configuration = configuration;
    }

    [Fact]
    public void GetRecipes()
    {

    }
}