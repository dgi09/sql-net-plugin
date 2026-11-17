using System.Text.Json;
using Microsoft.Extensions.DependencyInjection;
using Project.Data;

var cs = "Host=localhost;Port=5434;Database=sqlc_test;Username=postgres;Password=postgres";

IServiceCollection services = new ServiceCollection();
services.AddDb(cs);

IServiceProvider provider = services.BuildServiceProvider();

var db = provider.GetRequiredService<Db>();
// await Operation(db);
//

var query = new Query(db);
// await query.CreateAuthor("ivan", "bio", ["some tag"]);
var authors = await query.ListAuthors(minId: null, minCreatedAt: null, sortBy: "created_at");

Console.WriteLine(JsonSerializer.Serialize(authors));
