namespace Tours.Models;

public enum TransportType { car, bike, walk};
public class TourDuration
{
    public int Id { get; set; }
    public double Duration { get; set; }
    public TransportType TransportType { get; set; }

    public TourDuration(double duration, TransportType transportType)
    {
        Duration = duration;
        TransportType = transportType;
    }
}
