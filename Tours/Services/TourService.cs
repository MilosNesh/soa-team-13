using Tours.Models;
using Tours.Repositorues;
using FluentResults;

namespace Tours.Services;

public class TourService: ITourService
{
    private readonly ITourRepository _tourRepository;

    public TourService(ITourRepository repository)
    {
        _tourRepository = repository;
    }

    public Result<Tour> Create(Tour tour)
    {
        try
        {
            tour.Status = TourStatus.Draft;
            var savedTour = _tourRepository.Create(tour);
            return savedTour;
        }
        catch (Exception e)
        {
            return Result.Fail(new Error("Invalid data supplied.").WithMetadata("code", 400)).WithError(e.Message);
        }
    }

    public Result<List<Tour>> GetAll() 
    {
        return _tourRepository.GetAll();
    }

}
