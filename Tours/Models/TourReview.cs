namespace Tours.Models
{
    public class TourReview
    {
        public int Id { get; set; }
        public int Rating { get; set; }
        public string? Comment { get; set; }
        public DateTime TourDate { get; set; }
        public DateTime CreationDate { get; set; }
        public int TouristId { get; set; }
        public int TourId { get; set; }
        public TourReview(int rating, string comment, DateTime tourDate, DateTime creationDate, int touristId, int tourId)
        {
            Rating = rating;
            Comment = comment;
            TourDate = tourDate;
            CreationDate = creationDate;
            Validate();
            TouristId = touristId;
            TourId = tourId;
        }

        private void Validate()
        {
            if (Rating < 1 || Rating > 5)
                throw new ArgumentException("Rating must be between 1 and 5.");

            if (string.IsNullOrWhiteSpace(Comment))
                throw new ArgumentException("Comment cannot be empty. ");

            if (TourDate > DateTime.Now)
                throw new ArgumentException("Tour date cannot be in the future.");

            if (CreationDate < TourDate)
                throw new ArgumentException("Creation date cannot be before the tour date.");

        }

        public override string ToString()
        {
            return $"Rating: {Rating}, Comment: {Comment ?? "None"}, TourDate: {TourDate:g}, " +
                   $"CreationDate: {CreationDate:g}";
        }
    }
}
