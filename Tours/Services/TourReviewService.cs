using FluentResults;
using Tours.Models;
using Tours.Repositorues;

namespace Tours.Services
{
    public class TourReviewService : ITourReviewService
    {
        private readonly ITourReviewRepository _tourReviewRepository;
        public TourReviewService(ITourReviewRepository repository) 
        {
            _tourReviewRepository = repository;
        }
        public Result<List<TourReview>> GetAllByTourId(int tourId)
        {
            var reviews = _tourReviewRepository.GetByTourId(tourId);

            if (reviews == null || !reviews.Any())
                return Result.Fail<List<TourReview>>("Nema recenzija za ovu turu.");

            return Result.Ok(reviews);
        }

        public Result<TourReview> Create(TourReview tourReview)
        {
            try
            {
                var savedTourReview = _tourReviewRepository.Create(tourReview);
                return savedTourReview;
            }
            catch (Exception e)
            {
                return Result.Fail(new Error("Invalid data supplied.").WithMetadata("code", 400)).WithError(e.Message);
            }
        }

    }
}
