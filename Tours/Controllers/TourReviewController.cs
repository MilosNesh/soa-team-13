using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using Tours.Models;
using Tours.Services;

namespace Tours.Controllers
{
    [ApiController]
    [Route("tourReview")]
    public class TourReviewController : ControllerBase
    {
        private readonly ITourReviewService _tourReviewService;

        public TourReviewController(ITourReviewService tourReviewService)
        {
            _tourReviewService = tourReviewService;
        }

        [HttpGet("tour/{tourId:int}")]
        public ActionResult GetAll(int tourId)
        {
            var result = _tourReviewService.GetAllByTourId(tourId);

            if (result == null || !result.Value.Any())
            {
                return NotFound(new { message = "Nema recenzija za ovu turu." });
            }

            return Ok(result.Value);
        }


        [HttpPost]
        public ActionResult<TourReview> Create([FromBody] TourReview tourReview)
        {
            var result = _tourReviewService.Create(tourReview);
            return result.IsSuccess ? Ok(result.Value) : BadRequest(result.Errors);
        }
    }
}
