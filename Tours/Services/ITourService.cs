using FluentResults;
using Tours.Models;

namespace Tours.Services;

public interface ITourService
{
    public Result<Tour> Create(Tour tour);
    public Result<Tour> Update(Tour tour);
    public Result<List<Tour>> GetAll();
    public Result<List<Tour>> GetByAuthor(string authorId);
    public Result<Tour> GetById(int id);
    public Result<Tour> Publish(Tour tour);
    public Result<Tour> Archive(Tour tour);
    public Result<Tour> Reactivate(Tour tour);
    public Result<List<Tour>> GetPublished();
}
