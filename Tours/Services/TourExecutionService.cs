using FluentResults;
using Microsoft.IdentityModel.Tokens;
using Tours.Dtos;
using Tours.Models;
using Tours.Repositorues;

namespace Tours.Services
{
    public class TourExecutionService : ITourExecutionService
    {
        private readonly ITourExecutionRepository _tourExecutionRepository;
        private readonly ITourService _tourService;
        public TourExecutionService(ITourExecutionRepository tourExecutionRepository, ITourService tourService)
        {
            _tourExecutionRepository = tourExecutionRepository;
            _tourService = tourService;
        }

        public Result<List<TourExecution>> GetAll()
        {
            var result = _tourExecutionRepository.GetAll();
            if (result == null)
            {

                return Result.Fail("Tour executions not found");
            }
            return result;
        }

        public Result<TourExecution> Create(CreateTourExecutionDto dto)
        {
            try
            {
                var existingSession = _tourExecutionRepository.GetByUserId(dto.UserId);
                if (existingSession != null && existingSession.ExecutionStatus == ExecutionStatus.Active)
                {
                    return Result.Fail("This action is forbidden because execution is active");
                }

                var tourResult = _tourService.GetById(dto.TourId);
                if (tourResult.IsFailed)
                    return Result.Fail("Tour not found");

                var tour = tourResult.Value;
                var tourExecution = new TourExecution(tour.Id, dto.UserId, tour.Length);

                if (tour.KeyPoints == null)
                {
                    // dodaj log ili throw ovde da vidiš
                    Console.WriteLine("Tour.KeyPoints is null!");
                }

                tourExecution.AddKeyPointStatuses(tour.KeyPoints?.Where(kp => kp != null).ToList() ?? new List<KeyPoint>());

                _tourExecutionRepository.Create(tourExecution);
                return tourExecution;
            }
            catch (Exception ex)
            {
                // loguj i throw dalje
                Console.WriteLine(ex);
                throw;
            }
        }

        public Result<TourExecution> Update(int tourExecutionId, double longitude, double latitude)
        {
            var tourExecution = _tourExecutionRepository.Get(tourExecutionId);
            foreach (KeyPointStatus checkpointStatus in tourExecution.KeyPointStatus)
            {
                    tourExecution.UpdateLocation(longitude, latitude, checkpointStatus);
            }
            _tourExecutionRepository.Update(tourExecution);

            return tourExecution;
        }
        public Result<TourExecution> Get(int id)
        {
            var tourExecution = _tourExecutionRepository.Get(id);
            return tourExecution;
        }
        public async Task<int> GetTourCompletionPercentageAsync(int tourExecutionId)
        {
            var tourExecution = await _tourExecutionRepository.GetTourExecutionByIdAsync(tourExecutionId);
            return tourExecution?.GetTourCompletionPercentage() ?? 0;
        }

        public Result<TourExecution> GetSessionsByUserId(int userId)
        {
            var session = _tourExecutionRepository.GetByUserId(userId);
            if (session != null)
                return session;
            return Result.Fail("There is no active session to show");
        }

        public Result<TourExecution> CompleteSession(int userId)
        {
            var tourExecution = _tourExecutionRepository.GetByUserId(userId);

            if (tourExecution == null)
            {
                throw new Exception("Tour execution not found.");
            }

            if (tourExecution.KeyPointStatus.All(cs => cs.IsCompleted()))
            {
                tourExecution.CompleteSession();
                _tourExecutionRepository.Update(tourExecution);
                return tourExecution;
            }
            return Result.Fail("You cant complete this session");
        }

        public Result<TourExecution> AbandonSession(int userId)
        {
            var tourExecution = _tourExecutionRepository.GetByUserId(userId);

            if (tourExecution == null)
            {
                return Result.Fail("Tour execution not found.");
            }

            tourExecution.AbandonSession();
            _tourExecutionRepository.Update(tourExecution);
            return tourExecution;
        }
        public Result<List<int>> GetCompletedToursByTourist(int id)
        {
            try
            {
                var tourIds = _tourExecutionRepository.GetCompletedToursByTourist(id);
                return tourIds;
            }
            catch (Exception e)
            {
                return Result.Fail(e.ToString());
            }
        }

    }
}
