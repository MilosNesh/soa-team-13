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


var app = builder.Build();

// Configure the HTTP request pipeline.

//app.UseHttpsRedirection();

app.ApplyMigrations();

app.UseAuthorization();

app.MapControllers();

app.Run();
