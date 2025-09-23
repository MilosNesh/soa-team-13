using Microsoft.EntityFrameworkCore;
using Tours.Models;

namespace Tours.Repositorues
{
    public class TourExecutionRepository : ITourExecutionRepository
    {
        private readonly ToursContext _context;

        public TourExecutionRepository(ToursContext context)
        {
            _context = context;
        }

        public List<TourExecution> GetAll()
        {
            return _context.TourExecutions.ToList();
        }
        public new TourExecution Get(int id)
        {
            return _context.TourExecutions.Where(te => te.Id == id)
                .Include(te => te.KeyPointStatus)
                .ThenInclude(cs => cs.KeyPoint).FirstOrDefault() ?? throw new Exception("Id not found");
        }
        public new TourExecution Update(TourExecution tourExecution)
        {
            _context.Entry(tourExecution).State = EntityState.Modified;
            _context.SaveChanges();
            return tourExecution;
        }
        public async Task<TourExecution> GetTourExecutionByIdAsync(int tourExecutionId)
        {
            return await _context.TourExecutions
                                 .Include(te => te.KeyPointStatus)
                                 .FirstOrDefaultAsync(te => te.Id == tourExecutionId);
        }

        public TourExecution GetByUserId(int userId)
        {
            return _context.TourExecutions.Include(te => te.KeyPointStatus)
                    .ThenInclude(cs => cs.KeyPoint).FirstOrDefault(te => te.UserId == userId && te.ExecutionStatus == 0);
        }

        public List<int> GetCompletedToursByTourist(int userId)
        {
            var result = _context.TourExecutions
         .Where(te => te.UserId == userId && te.ExecutionStatus == ExecutionStatus.Completed)
         .Select(te => te.TourId)
         .ToList();

            return result;
        }

        public TourExecution Create(TourExecution tourExecution)
        {
            try
            {
                tourExecution.StartTime = DateTime.SpecifyKind(tourExecution.StartTime, DateTimeKind.Utc);
                tourExecution.LastActivity = DateTime.SpecifyKind(tourExecution.LastActivity, DateTimeKind.Utc);

                _context.TourExecutions.Add(tourExecution);
                _context.SaveChanges();

                return tourExecution;
            }
            catch (Exception)
            {
                throw;
            }
        }

    }
}
