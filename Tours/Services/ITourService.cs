using FluentResults;
using Tours.Models;

namespace Tours.Services;

public interface ITourService
{
    public Result<Tour> Create(Tour tour);
    public Result<List<Tour>> GetAll();

}
