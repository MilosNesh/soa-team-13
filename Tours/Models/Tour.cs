using System.Security.Cryptography;

namespace Tours.Models;

public enum TourStatus
{
    Draft = 0,
    Published = 1,
    Archived = 2,
    Closed = 3
}

public enum TourDifficulty
{
    Easy = 0,
    Medium = 1,
    Hard = 2,
    Hell = 3
}
public class Tour 
{
    public int Id { get; set; }
    public string Name { get;  set; }
    public TourDifficulty Difficulty { get;  set; }
    public string Description { get;  set; }
    public double Cost { get;  set; }
    public TourStatus Status { get;  set; }
    public string Tags { get;  set; }
    public List<KeyPoint> KeyPoints { get; set; } = new List<KeyPoint>();
    public List<TourDuration> Durations { get; set; } = new List<TourDuration>();
    public double Length { get;  set; }
    public string AuthorId { get;  set; }
    public DateTime? PublishTime { get;  set; } = null;
    public DateTime? ArchiveTime { get;  set; } = null;
    public string Image { get;  set; }
    public List<TourReview> Reviews { get; set; }

    public Tour(string name, TourDifficulty difficulty, string description, double cost, TourStatus status, string tags, double length, string authorId, string image)
    {
        if (string.IsNullOrWhiteSpace(name)) throw new ArgumentException("Invalid Name.");
        Name = name;
        Difficulty = difficulty;
        Description = description;
        Cost = cost;
        Status = status;
        Tags = tags;
        Length = length;
        AuthorId = authorId;
        Image = image;
        Reviews = new List<TourReview>();
    }

    public Tour()
    {
    }

    //public void AddKeyPoint(KeyPoint keyPoint)
    //{
    //    if (keyPoint == null) throw new ArgumentNullException(nameof(keyPoint));
    //    KeyPoints.Add(keyPoint);
    //}

    public Tour Publish()
    {
        if (!CanPublish())
            throw new ArgumentException("Tura nije ispunila uslove za objavljivanje.");
        Status = TourStatus.Published;
        PublishTime = DateTime.UtcNow;
        return this;
    }

    public bool CanPublish()
    {
        return !string.IsNullOrEmpty(Name) && !string.IsNullOrEmpty(Description) && !string.IsNullOrEmpty(Description) && !string.IsNullOrEmpty(Tags)
               && KeyPoints.Count >= 2 && Durations.Count >= 1;
    }

    public Tour Archive()
    {
        if (!CanArchive())
            throw new ArgumentException("Nije moguce arhivirati turu jer");
        Status = TourStatus.Archived;
        ArchiveTime = DateTime.UtcNow;
        return this;
    }

    public Tour UpdateTourLength(double length)
    {
        Length = length;
        return this;
    }

    public bool CanArchive()
    {
        return Status == TourStatus.Published;
    }

    public Tour Reactivate()
    {
        if (!CanReactivate())
            throw new ArgumentException("Nije moguce re-aktivirati ovu turu jer nije arhivirana");
        Status = TourStatus.Published;
        return this;
    }
    public bool CanReactivate()
    {
        return Status == TourStatus.Archived;
    }

    public void CloseTour()
    {
        Status = TourStatus.Closed;
    }

    public Tour Preview()
    {
        var keyPoints = new List<KeyPoint>();
        keyPoints.Add(KeyPoints.FirstOrDefault());
        return new Tour
        {
            Id = Id,
            Name = Name,
            Description = Description,
            Tags = Tags,
            Image = Image,
            Cost = Cost,
            KeyPoints = keyPoints
        };
    }
}
