import 'clothing_item.dart';

class Outfit {
  final String id;
  final String userId;
  final String name;
  final String? description;
  final String? style;
  final List<String> occasion;
  final List<String> season;
  final String? weather;
  final double? colorHarmonyScore;
  final double? styleCoherenceScore;
  final double? overallScore;
  final int? rating;
  final bool worn;
  final bool favorite;
  final String? notes;
  final DateTime? createdAt;
  final DateTime? updatedAt;
  final List<ClothingItem> clothingItems;

  Outfit({
    required this.id,
    required this.userId,
    required this.name,
    this.description,
    this.style,
    this.occasion = const [],
    this.season = const [],
    this.weather,
    this.colorHarmonyScore,
    this.styleCoherenceScore,
    this.overallScore,
    this.rating,
    this.worn = false,
    this.favorite = false,
    this.notes,
    this.createdAt,
    this.updatedAt,
    this.clothingItems = const [],
  });

  factory Outfit.fromJson(Map<String, dynamic> json) {
    return Outfit(
      id: json['id'] ?? '',
      userId: json['user_id'] ?? '',
      name: json['name'] ?? '',
      description: json['description'],
      style: json['style'],
      occasion: List<String>.from(json['occasion'] ?? []),
      season: List<String>.from(json['season'] ?? []),
      weather: json['weather'],
      colorHarmonyScore: json['color_harmony_score']?.toDouble(),
      styleCoherenceScore: json['style_coherence_score']?.toDouble(),
      overallScore: json['overall_score']?.toDouble(),
      rating: json['rating'],
      worn: json['worn'] ?? false,
      favorite: json['favorite'] ?? false,
      notes: json['notes'],
      createdAt: json['created_at'] != null 
          ? DateTime.tryParse(json['created_at']) 
          : null,
      updatedAt: json['updated_at'] != null 
          ? DateTime.tryParse(json['updated_at']) 
          : null,
      clothingItems: (json['clothing_items'] as List?)
          ?.map((e) => ClothingItem.fromJson(e))
          .toList() ?? [],
    );
  }
}
