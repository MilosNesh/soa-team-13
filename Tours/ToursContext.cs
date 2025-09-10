using Microsoft.EntityFrameworkCore;
using Tours.Models;

namespace Tours;

public class ToursContext : DbContext
{
    public ToursContext(DbContextOptions<ToursContext> options) : base(options)
    {
    }
    public DbSet<Tour> Tours { get; set; }
    public DbSet<KeyPoint> KeyPoints { get; set; }

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        //modelBuilder.HasDefaultShames("tours");
        modelBuilder.Entity<Tour>().
            HasMany(t => t.KeyPoints).
            WithOne();

    }
}
