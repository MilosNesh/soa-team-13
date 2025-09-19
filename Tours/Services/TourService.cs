using Tours.Models;
using Tours.Repositorues;
using FluentResults;

namespace Tours.Services;

public class TourService : ITourService
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

    public Result<Tour> Update(Tour tour)
    {
        try
        {
            var savedTour = _tourRepository.Update(tour);
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

    public Result<List<Tour>> GetByAuthor(string authorId)
    {
        return _tourRepository.GetByAuthorId(authorId);
    }

    public Result<Tour> GetById(int id)
    {
        return _tourRepository.GetById(id);
    }

    public Result<Tour> Publish(Tour tour)
    {
        try
        {
            Tour t = tour.Publish();
            return _tourRepository.Update(t);
        }
        catch (ArgumentException e)
        {
            return Result.Fail(new Error("Invalid data supplied.")
                .WithMetadata("code", 400)).WithError(e.Message);
        }
    }

    public Result<Tour> Archive(Tour tour)
    {
        try
        {
            Tour t = tour.Archive();
            return _tourRepository.Update(t);
        }
        catch (ArgumentException e)
        {
            return Result.Fail(new Error("Invalid data supplied.")
                .WithMetadata("code", 400)).WithError(e.Message);
        }
    }

    public Result<Tour> Reactivate(Tour tour)
    {
        try
        {
            Tour t = tour.Reactivate();
            return _tourRepository.Update(t);
        }
        catch (ArgumentException e)
        {
            return Result.Fail(new Error("Invalid data supplied.")
                .WithMetadata("code", 400)).WithError(e.Message);
        }
    }

    public Result<List<Tour>> GetPublished()
    {
        return _tourRepository.GetPublished();
    }
}
