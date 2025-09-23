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
    public DbSet<TourReview> TourReviews { get; set; }
    public DbSet<TourExecution> TourExecutions { get; set; }
    public DbSet<KeyPointStatus> KeyPointStatuses { get; set; }

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        //modelBuilder.HasDefaultShames("tours");
        modelBuilder.Entity<Tour>().
            HasMany(t => t.KeyPoints).
            WithOne();
        modelBuilder.Entity<Tour>().HasMany(tr => tr.Reviews).WithOne().HasForeignKey(r => r.TourId);

        modelBuilder.Entity<Tour>().HasMany(t => t.Durations).WithOne();

        modelBuilder.Entity<TourExecution>().HasMany(te => te.KeyPointStatus).WithOne().OnDelete(DeleteBehavior.Cascade);
        modelBuilder.Entity<KeyPointStatus>()
            .HasOne(cs => cs.KeyPoint)
            .WithMany()
            .HasForeignKey(cs => cs.KeyPointId);
    }
}
