namespace Tours.Models
{
    public class TourExecution
    {
        public int Id { get; set; }
        public int TourId { get; set; }
        public int UserId { get; set; }
        public double TourRange { get; set; }
        public DateTime StartTime { get; set; }
        public DateTime? EndTime { get; set; }
        public DateTime LastActivity { get; set; }
        public ExecutionStatus ExecutionStatus { get; set; }
        public ICollection<KeyPointStatus> KeyPointStatus { get; set; } = new List<KeyPointStatus>();
        public TourExecution(int tourId, int userId, double tourRange)
        {
            TourId = tourId;
            UserId = userId;
            TourRange = tourRange;
            StartTime = DateTime.UtcNow;
            LastActivity = DateTime.UtcNow;
            ExecutionStatus = ExecutionStatus.Active;
        }

        public void CompleteSession()
        {
            if (ExecutionStatus == ExecutionStatus.Active && KeyPointStatus.All(cs => cs.IsCompleted()))
            {
                ExecutionStatus = ExecutionStatus.Completed;
                EndTime = DateTime.UtcNow;
                LastActivity = DateTime.UtcNow;
            }
        }

        public void AbandonSession()
        {
            if (ExecutionStatus == ExecutionStatus.Active)
            {
                ExecutionStatus = ExecutionStatus.Abandoned;
                EndTime = DateTime.UtcNow;
                LastActivity = DateTime.UtcNow;
            }
        }
        public void UpdateLocation(double longitude, double latitude, KeyPointStatus keyPointStatus)
        {
            LastActivity = DateTime.UtcNow;
            if (!keyPointStatus.IsCompleted() && keyPointStatus.IsTouristNear(latitude, longitude))
            {
                keyPointStatus.MarkAsCompleted();
            }
        }

        public void AddKeyPointStatuses(List<KeyPoint> keyPoints)
        {
            foreach (KeyPoint keyPoint in keyPoints)
            {
                KeyPointStatus.Add(new KeyPointStatus(keyPoint.Id));
            }
        }

        public int GetTourCompletionPercentage()
        {
            if (KeyPointStatus == null || !KeyPointStatus.Any())
                return 0;

            int completedKeyPoints = KeyPointStatus.Count(c => c.IsCompleted());
            double completionPercentage = (double)completedKeyPoints / KeyPointStatus.Count * 100;

            return (int)Math.Round(completionPercentage);
        }
    }

    public enum ExecutionStatus
    {
        Active,
        Completed,
        Abandoned
    }
}
