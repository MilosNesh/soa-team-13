using Tours.Models;

namespace Tours.Repositorues;

public interface ITourRepository
{
    public Tour Create(Tour tour);
    public Tour Get(int id);
    public List<Tour> GetAll();
    public List<Tour> GetByAuthorId(string id);
    public Tour Update(Tour tour);
}
