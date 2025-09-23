using FluentResults;
using Tours.Dtos;
using Tours.Models;

namespace Tours.Services
{
    public interface ITourExecutionService
    {
        Result<List<TourExecution>> GetAll();
        Result<TourExecution> Update(int tourExecutionId, double longitude, double latitude);
        Result<TourExecution> Create(CreateTourExecutionDto tourExecution);
        Result<TourExecution> GetSessionsByUserId(int userId);
        Result<TourExecution> CompleteSession(int userId);
        Result<TourExecution> AbandonSession(int userId);
        Result<TourExecution> Get(int id);
        Task<int> GetTourCompletionPercentageAsync(int tourExecutionId);
        Result<List<int>> GetCompletedToursByTourist(int id);
    }
}
