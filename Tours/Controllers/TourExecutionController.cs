using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using Tours.Dtos;
using Tours.Models;
using Tours.Services;

namespace Tours.Controllers
{
    [ApiController]
    [Route("tourExecutions/")]
    public class TourExecutionController : ControllerBase
    {
        private readonly ITourExecutionService _tourExecutionService;

        public TourExecutionController(ITourExecutionService tourExecutionService)
        {
            _tourExecutionService = tourExecutionService;
        }

        [HttpGet]
        public ActionResult<List<TourExecution>> GetAll()
        {
            var result = _tourExecutionService.GetAll();
            return result.IsSuccess ? Ok(result.Value) : NotFound(result.Errors);

        }

        [HttpGet("{id:int}")]
        public ActionResult<TourExecution> GetById(int id)
        {
            var result = _tourExecutionService.Get(id);
            return result.IsSuccess ? Ok(result.Value) : BadRequest(result.Errors);

        }
        [HttpGet("completed/{id:int}")]
        public ActionResult<List<long>> GetCompletedToursByTourist(int id)
        {
            var result = _tourExecutionService.GetCompletedToursByTourist(id);
            return result.IsSuccess ? Ok(result.Value) : BadRequest(result.Errors);

        }

        [HttpPost]
        public ActionResult<TourExecution> Create([FromBody] CreateTourExecutionDto tourExecution)
        {
            try
            {
                var result = _tourExecutionService.Create(tourExecution);
                return result.IsSuccess ? Ok(result.Value) : BadRequest(result.Errors);
            }
            catch (Exception ex)
            {
                return StatusCode(500, new { error = ex.Message, stack = ex.StackTrace });
            }

        }
        [HttpPut("update/{tourExecutionId:int}/{longitude:double}/{latitude:double}")]
        public ActionResult<TourExecution> Update([FromRoute] int tourExecutionId, double longitude, double latitude)
        {
            var result = _tourExecutionService.Update(tourExecutionId, longitude, latitude);
            return result.IsSuccess ? Ok(result.Value) : BadRequest(result.Errors);

        }

        [HttpGet("{tourExecutionId}/completion-percentage")]
        public async Task<ActionResult<int>> GetTourCompletionPercentage(int tourExecutionId)
        {
            int completionPercentage = await _tourExecutionService.GetTourCompletionPercentageAsync(tourExecutionId);
            return Ok(completionPercentage);
        }

        [HttpGet("getByUser/{userId:int}")]
        public ActionResult<TourExecution> GetByUser(int userId)
        {
            var result = _tourExecutionService.GetSessionsByUserId(userId);
            return result.IsSuccess ? Ok(result.Value) : BadRequest(result.Errors);
        }

        // za swagger

        [HttpPut("complete/{userId}")]
        public ActionResult<TourExecution> CompleteSession(int userId)
        {
            var result = _tourExecutionService.CompleteSession(userId);
            return result.IsSuccess ? Ok(result.Value) : BadRequest(result.Errors);
        }

        [HttpPut("abandon/{userId}")]
        public ActionResult<TourExecution> AbandonSession(int userId)
        {
            var result = _tourExecutionService.AbandonSession(userId);
            return result.IsSuccess ? Ok(result.Value) : BadRequest(result.Errors);

        }
    }
}
