using Microsoft.AspNetCore.Mvc;
using Tours.Models;
using Tours.Services;

namespace Tours.Controllers
{

    [ApiController]
    [Route("key-points/")]
    public class KeyPointController : ControllerBase
    {
       private readonly IKeyPointService _keyPointService;

        public KeyPointController(IKeyPointService keyPointService)
        {
            _keyPointService = keyPointService;
        }

        [HttpPost]
        public ActionResult<KeyPoint> Create([FromBody] KeyPoint keyPoint)
        {
            var result = _keyPointService.Create(keyPoint);
            return result.IsSuccess ? Ok(result.Value) : BadRequest(result.Errors);
        }

        [HttpGet("{id:int}")]
        public ActionResult<KeyPoint> GetById(int id)
        {
            var result = _keyPointService.Get(id);
            return result.IsSuccess ? Ok(result.Value) : NotFound(result.Errors);
        }
    }
}
