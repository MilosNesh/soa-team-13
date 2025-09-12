using FluentResults;
using Tours.Models;
using Tours.Repositorues;

namespace Tours.Services
{
    public class KeyPointService: IKeyPointService
    {
        private readonly IKeyPointRepository _keyPointRepository;

        public KeyPointService(IKeyPointRepository repository)
        {
            _keyPointRepository = repository;
        }

        public Result<KeyPoint> Create(KeyPoint keyPoint)
        {
            try
            {
                var savedKeyPoint = _keyPointRepository.Create(keyPoint);
                return savedKeyPoint;
            }
            catch (Exception e)
            {
                return Result.Fail(new Error("Invalid data supplied.").WithMetadata("code", 400)).WithError(e.Message);
            }
        }

        public Result<KeyPoint> Get(int id)
        {
            return _keyPointRepository.Get(id);
        }

        public Result Delete(int id)
        {
            return _keyPointRepository.Delete(id);
        }
    }
}
