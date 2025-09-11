using Tours.Models;

namespace Tours.Repositorues
{
    public interface ITourReviewRepository
    {
        public List<TourReview> GetByTourId(int id);
        public TourReview Create(TourReview tourReview);
    }
}
