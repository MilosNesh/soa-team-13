using Microsoft.EntityFrameworkCore;
using Tours;
using Tours.Repositorues;
using Tours.Services;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddDbContext<ToursContext>(options =>
    options.UseNpgsql(builder.Configuration.GetConnectionString("ToursDatabase")));

builder.Services.AddScoped<ITourRepository, TourRepository>();
builder.Services.AddScoped<ITourService, TourService>(); 
builder.Services.AddControllers()
        .AddJsonOptions(options => options.JsonSerializerOptions.PropertyNamingPolicy = System.Text.Json.JsonNamingPolicy.CamelCase);
builder.Services.AddScoped<ITourReviewService, TourReviewService>();
builder.Services.AddScoped<ITourReviewRepository, TourReviewRepository>();
builder.Services.AddScoped<IKeyPointRepository, KeyPointRepository>();
builder.Services.AddScoped<IKeyPointService, KeyPointService>();
builder.Services.AddScoped<ITourExecutionService, TourExecutionService>();
builder.Services.AddScoped<ITourExecutionRepository, TourExecutionRepository>();

builder.Services.AddCors(options =>
{
    options.AddPolicy("DevCors", policy =>
    {
        policy
            .WithOrigins("http://localhost:4200")
            .AllowAnyHeader()
            .AllowAnyMethod();
    });
});

var app = builder.Build();

// Configure the HTTP request pipeline.

//app.UseHttpsRedirection();

app.ApplyMigrations();

app.UseAuthorization();

app.UseCors("DevCors");

app.MapControllers();

app.Run();
