namespace Tours.Models;

public class KeyPoint
{
    public int Id { get; set; }
    public string Name { get; init; }
    public string Description { get; init; }
    public string Image { get; init; }
    public float Latitude { get; init; }
    public float Longitude { get; init; }

    public KeyPoint(string name, string description, string image, float latitude, float longitude)
    {
        if (string.IsNullOrWhiteSpace(name)) throw new ArgumentNullException("Invalid name");
        if (string.IsNullOrWhiteSpace(description)) throw new ArgumentNullException("Invalid description");

        Name = name;
        Description = description;
        Image = image;
        Latitude = latitude;
        Longitude = longitude;
     
    }
    public KeyPoint(KeyPoint keyPoint)
    {
        Name = keyPoint.Name;
        Description = keyPoint.Description;
        Image = keyPoint.Image;
        Latitude = keyPoint.Latitude;
        Longitude = keyPoint.Longitude;
    }

}
