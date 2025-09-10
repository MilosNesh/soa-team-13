using FluentResults;
using Tours.Models;

namespace Tours.Services
{
    public interface ITourReviewService
    {
        Result<List<TourReview>> GetAllByTourId(int tourId);

        Result<TourReview> Create(TourReview tourReview);

    }
}
