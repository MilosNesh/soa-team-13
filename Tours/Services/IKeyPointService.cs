using FluentResults;
using Tours.Models;

namespace Tours.Services
{
    public interface IKeyPointService
    {
        public Result<KeyPoint> Create(KeyPoint keyPoint);
        public Result<KeyPoint> Get(int id);
        public Result Delete(int id);

    }
}
