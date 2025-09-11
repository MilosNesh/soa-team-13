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
        return _context.Tours.Include(t => t.KeyPoints).FirstOrDefault(t => t.Id == id);
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
        var list = (List<Tour>) _context.Tours.Include(t => t.KeyPoints).Where(t => t.AuthorId == id);
        return list;
    }

    public List<Tour> GetAll()
    {
        var list = _context.Tours.Include(t => t.KeyPoints).ToList();
        return list;
    }

    public Tour GetById(int id)
    {
        return _context.Tours
            .AsNoTracking()
            .Include(t => t.KeyPoints)
            .FirstOrDefault(t => t.Id == id);
    }

}
