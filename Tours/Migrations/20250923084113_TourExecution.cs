using System;
using Microsoft.EntityFrameworkCore.Migrations;
using Npgsql.EntityFrameworkCore.PostgreSQL.Metadata;

#nullable disable

namespace Tours.Migrations
{
    /// <inheritdoc />
    public partial class TourExecution : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.CreateTable(
                name: "TourExecutions",
                columns: table => new
                {
                    Id = table.Column<int>(type: "integer", nullable: false)
                        .Annotation("Npgsql:ValueGenerationStrategy", NpgsqlValueGenerationStrategy.IdentityByDefaultColumn),
                    TourId = table.Column<int>(type: "integer", nullable: false),
                    UserId = table.Column<int>(type: "integer", nullable: false),
                    TourRange = table.Column<double>(type: "double precision", nullable: false),
                    StartTime = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    EndTime = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    LastActivity = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    ExecutionStatus = table.Column<int>(type: "integer", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_TourExecutions", x => x.Id);
                });

            migrationBuilder.CreateTable(
                name: "KeyPointStatuses",
                columns: table => new
                {
                    Id = table.Column<int>(type: "integer", nullable: false)
                        .Annotation("Npgsql:ValueGenerationStrategy", NpgsqlValueGenerationStrategy.IdentityByDefaultColumn),
                    KeyPointId = table.Column<int>(type: "integer", nullable: false),
                    TourExecutionId = table.Column<long>(type: "bigint", nullable: false),
                    CompletionTime = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    TourExecutionId1 = table.Column<int>(type: "integer", nullable: true)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_KeyPointStatuses", x => x.Id);
                    table.ForeignKey(
                        name: "FK_KeyPointStatuses_KeyPoints_KeyPointId",
                        column: x => x.KeyPointId,
                        principalTable: "KeyPoints",
                        principalColumn: "Id",
                        onDelete: ReferentialAction.Cascade);
                    table.ForeignKey(
                        name: "FK_KeyPointStatuses_TourExecutions_TourExecutionId1",
                        column: x => x.TourExecutionId1,
                        principalTable: "TourExecutions",
                        principalColumn: "Id",
                        onDelete: ReferentialAction.Cascade);
                });

            migrationBuilder.CreateIndex(
                name: "IX_KeyPointStatuses_KeyPointId",
                table: "KeyPointStatuses",
                column: "KeyPointId");

            migrationBuilder.CreateIndex(
                name: "IX_KeyPointStatuses_TourExecutionId1",
                table: "KeyPointStatuses",
                column: "TourExecutionId1");
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropTable(
                name: "KeyPointStatuses");

            migrationBuilder.DropTable(
                name: "TourExecutions");
        }
    }
}
