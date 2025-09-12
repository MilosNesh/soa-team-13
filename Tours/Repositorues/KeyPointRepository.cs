using FluentResults;
using Microsoft.EntityFrameworkCore;
using Tours.Models;

namespace Tours.Repositorues
{
    public class KeyPointRepository: IKeyPointRepository
    {
        private readonly ToursContext _context;

        public KeyPointRepository(ToursContext context)
        {
            _context = context;
        }

        public KeyPoint Get(int id)
        {
            return _context.KeyPoints.FirstOrDefault(t => t.Id == id);
        }

        public KeyPoint Create(KeyPoint keyPoint)
        {
            try
            {
                _context.KeyPoints.Add(keyPoint);
                _context.SaveChanges();
                return keyPoint;
            }
            catch (Exception ex)
            {
                throw new Exception();
            }
        }

        public KeyPoint Update(KeyPoint keyPoint)
        {
            try
            {
                _context.Update(keyPoint);
                _context.SaveChanges();
                return keyPoint;
            }
            catch (Exception ex)
            {
                throw new Exception();
            }
        }
        public Result Delete(int id)
        {
            KeyPoint kp = _context.KeyPoints.FirstOrDefault(t => t.Id == id);
            _context.Remove(kp);
            var rows = _context.SaveChanges();
            return rows > 0 ? Result.Ok() : Result.Fail($"KeyPoint {id} not found");
        }
    }
}
