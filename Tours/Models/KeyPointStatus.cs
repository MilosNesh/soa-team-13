namespace Tours.Models
{
    public class KeyPointStatus
    {
        public int Id { get; set; }
        public int KeyPointId { get; private set; }
        public KeyPoint KeyPoint { get; private set; }
        public long TourExecutionId { get; private set; }
        public DateTime CompletionTime { get; private set; } = DateTime.MinValue;
        public KeyPointStatus(int keyPointId)
        {
            KeyPointId = keyPointId;
        }
        public bool IsCompleted()
        {
            return CompletionTime != DateTime.MinValue;
        }
        public bool IsTouristNear(double latitude, double longitude)
        {
            const double tolerance = 0.0018; // Tolerancija za blizinu (oko 11 metara)

            bool isNearLatitude = Math.Abs(KeyPoint.Latitude - latitude) <= tolerance;
            bool isNearLongitude = Math.Abs(KeyPoint.Longitude - longitude) <= tolerance;

            return isNearLatitude && isNearLongitude;
        }
        public void MarkAsCompleted()
        {
            CompletionTime = DateTime.UtcNow;
        }
    }
}
