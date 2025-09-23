using Tours.Models;

namespace Tours.Repositorues
{
    public interface ITourExecutionRepository
    {
        List<TourExecution> GetAll();
        TourExecution Get(int id);
        TourExecution Update(TourExecution tourExecution);
        Task<TourExecution> GetTourExecutionByIdAsync(int tourExecutionId);
        TourExecution GetByUserId(int userId);
        List<int> GetCompletedToursByTourist(int userId);
        TourExecution Create(TourExecution tourExecution);

    }
}