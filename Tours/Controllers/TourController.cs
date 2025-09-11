using FluentResults;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using Tours.Models;
using Tours.Services;

namespace Tours.Controllers;

[ApiController]
[Route("tours/")]
public class TourController : ControllerBase
{
    private readonly ITourService _tourService;

    public TourController(ITourService tourService)
    {
        _tourService = tourService;
    }

    [HttpPost]
    public ActionResult<Tour> Create([FromBody] Tour tour)
    {
        var result = _tourService.Create(tour);
        return result.IsSuccess ? Ok(result.Value) : BadRequest(result.Errors);
    }

    [HttpGet]
    public ActionResult<List<Tour>> GetAll()
    {
        var result = _tourService.GetAll();

        if (result.IsSuccess)
        {
            return Ok(result.Value);
        }

        return BadRequest(result.Errors);
    }

    [HttpGet("{id:int}")]
    public ActionResult<Tour> GetById(int id)
    {
        var result = _tourService.GetById(id);
        return result.IsSuccess ? Ok(result.Value) : BadRequest(result.Errors);
    }

    [HttpPut("{id:int}")]
    public ActionResult<Tour> Update([FromBody] Tour tour)
    {
        var result = _tourService.Update(tour);
        return result.IsSuccess ? Ok(result.Value) : BadRequest(result.Errors);
    }
}
