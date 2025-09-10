using Microsoft.EntityFrameworkCore;

namespace Tours;

public static class MigrationExtensions
{
    public static void ApplyMigrations(this IApplicationBuilder app)
    {
        using IServiceScope scope = app.ApplicationServices.CreateScope();

        using ToursContext toursContext = scope.ServiceProvider.GetRequiredService<ToursContext>();

        toursContext.Database.Migrate();
    }
}
