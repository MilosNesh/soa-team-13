using Microsoft.EntityFrameworkCore;
using Tours.Models;

namespace Tours.Repositorues
{
    public class TourReviewRepository : ITourReviewRepository
    {
        private readonly ToursContext _context;

        public TourReviewRepository(ToursContext context)
        {
            _context = context;
        }

        public TourReview Create(TourReview tourReview)
        {
            try
            {
                tourReview.TourDate = DateTime.SpecifyKind(tourReview.TourDate, DateTimeKind.Utc);
                tourReview.CreationDate = DateTime.SpecifyKind(tourReview.CreationDate, DateTimeKind.Utc);

                _context.TourReviews.Add(tourReview);
                _context.SaveChanges();
                return tourReview;
            }
            catch (Exception ex)
            {
                throw;
            }

        }

        public List<TourReview> GetByTourId(int id)
        {
            return _context.TourReviews
                   .Where(t => t.TourId == id)
                   .ToList();
        }
    }
}
