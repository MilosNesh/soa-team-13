using FluentResults;
using Tours.Models;

namespace Tours.Repositorues
{
    public interface IKeyPointRepository
    {
        public KeyPoint Get(int id);
        public KeyPoint Create(KeyPoint keyPoint);
        public KeyPoint Update(KeyPoint keyPoint);
        public Result Delete(int id);
    }
}
