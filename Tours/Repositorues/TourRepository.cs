using Microsoft.EntityFrameworkCore;
using Tours.Models;

namespace Tours.Repositorues;

public class TourRepository : ITourRepository
{
    private readonly ToursContext _context;

    public TourRepository(ToursContext context)
    {
        _context = context;
    }

    public Tour Get(int id)
    {
        return _context.Tours.Include(t => t.KeyPoints).Include(t => t.Durations).FirstOrDefault(t => t.Id == id);
    }

    public Tour Create(Tour tour)
    {
        try
        {
            _context.Tours.Add(tour);
            _context.SaveChanges();
            return tour;
        }
        catch (Exception ex)
        {
            throw new Exception();
        }
    }

    public Tour Update(Tour tour)
    {
        try
        {
            //DeleteDurations(tour);
            _context.Update(tour);
            _context.SaveChanges();
            return tour;
        }
        catch (Exception ex)
        {
            throw new Exception();
        }
    }

    public List<Tour> GetByAuthorId(string id)
    {
        var list = _context.Tours.Include(t => t.KeyPoints).Include(t => t.Durations).Where(t => t.AuthorId == id).ToList();
        return list;
    }

    public List<Tour> GetAll()
    {
        var list = _context.Tours.Include(t => t.KeyPoints).Include(t => t.Durations).ToList();
        return list;
    }

    public Tour GetById(int id)
    {
        return _context.Tours
            .AsNoTracking()
            .Include(t => t.KeyPoints).Include(t => t.Durations)
            .FirstOrDefault(t => t.Id == id);
    }

    public void DeleteDurations(Tour tour)
    {
        var existingTour = _context.Tours
            .Include(t => t.Durations)
            .FirstOrDefault(t => t.Id == tour.Id);

        existingTour.Durations.Clear();
        _context.SaveChanges();
    }

    public List<Tour> GetPublished()
    {
        var list = _context.Tours.Include(t => t.KeyPoints).Where(t => t.Status == TourStatus.Published).Select(t => t.Preview()).ToList();
        return list;
    }
}
